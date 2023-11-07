package info

import (
	"context"
	"fmt"

	hf "github.com/hobbyfarm/gargantua/pkg/apis/hobbyfarm.io/v1"

	// hfClientSet "github.com/hobbyfarm/gargantua/pkg/client/clientset/versioned/typed/hobbyfarm.io/v1"
	"os"
	"text/tabwriter"

	hfClientSet "github.com/hobbyfarm/gargantua/pkg/client/clientset/versioned"
	"github.com/sirupsen/logrus"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type SessionDetails struct {
	SessionVMMap map[string][]hf.VirtualMachine
}

func GetEmail(email string, Namespace string, hfc *hfClientSet.Clientset) (err error) {

	userid, err := getUser(email, Namespace, hfc)
	if err != nil {
		return err
	}

	sDetails, err := getUserAllocatedVMs(userid, Namespace, hfc)
	if err != nil {
		return err
	}

	return printReport([]SessionDetails{*sDetails})
}

func GetAccessCode(accesscode string, Namespace string, hfc *hfClientSet.Clientset, stats bool) (err error) {

	seList, err := hfc.HobbyfarmV1().ScheduledEvents(Namespace).List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		if apierrors.IsNotFound(err) {
			logrus.Infof("accesscode %s not found", accesscode)
			return nil
		} else {
			return err
		}
	}

	var seEvents []hf.ScheduledEvent

	for _, se := range seList.Items {
		if se.Spec.AccessCode == accesscode {
			seEvents = append(seEvents, se)
		}
	}

	if len(seEvents) > 1 {
		seNames := ""
		for _, se := range seEvents {
			seNames = seNames + " " + se.Spec.Name
		}
		return fmt.Errorf("more than one scheduledevent have the same access code: %s", seNames)
	}

	if len(seEvents) == 0 {
		logrus.Infof("No scheduled event with access code %s found", accesscode)
		return nil
	}

	logrus.Infof("scheduled event %s has accesscode %s", seEvents[0].Spec.Name, accesscode)
	userIDS, err := getAllUsers(accesscode, Namespace, hfc)
	if err != nil {
		return err
	}

	var sDetails []SessionDetails
	for _, user := range userIDS {
		s, err := getUserAllocatedVMs(user, Namespace, hfc)
		if err != nil {
			return err
		}
		sDetails = append(sDetails, *s)
	}

	if stats {
		return printStats(sDetails)
	}
	return printReport(sDetails)
}

func getUser(email string, Namespace string, hfc *hfClientSet.Clientset) (userid string, err error) {
	users, err := hfc.HobbyfarmV1().Users(Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			return userid, fmt.Errorf("no users found in cluster")
		} else {
			return userid, err
		}
	}

	for _, user := range users.Items {
		if user.Spec.Email == email {
			userid = user.Spec.Id
			break
		}
	}

	return userid, err
}

func getAllUsers(accesscode string, Namespace string, hfc *hfClientSet.Clientset) (userIDS []string, err error) {
	userList, err := hfc.HobbyfarmV1().Users(Namespace).List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		if apierrors.IsNotFound(err) {
			return userIDS, fmt.Errorf("no users found")
		} else {
			return userIDS, err
		}
	}

	for _, user := range userList.Items {
		if user.Spec.AccessCodes[0] == accesscode {
			// populate users with matching access code
			userIDS = append(userIDS, user.Spec.Id)
		}
	}

	return userIDS, nil
}

func getUserAllocatedVMs(userID string, Namespace string, hfc *hfClientSet.Clientset) (sDetails *SessionDetails, err error) {
	sessionList, err := hfc.HobbyfarmV1().Sessions(Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil, fmt.Errorf("no active sessions found for the user")
		} else {
			return nil, err
		}
	}

	tmpStore := make(map[string][]string)
	for _, session := range sessionList.Items {
		if session.Spec.UserId == userID {
			tmpStore[session.Spec.Id] = session.Spec.VmClaimSet
		}
	}

	sDetails = &SessionDetails{
		SessionVMMap: make(map[string][]hf.VirtualMachine),
	}

	for session, vmclaims := range tmpStore {
		var vmList []hf.VirtualMachine
		for _, vmc := range vmclaims {
			vml, err := getVirtualMachinesForVMC(vmc, Namespace, hfc)
			if err != nil {
				return nil, err
			}
			vmList = append(vmList, vml...)
		}
		sDetails.SessionVMMap[session] = vmList
	}

	return sDetails, nil
}

func getVirtualMachinesForVMC(vmc string, Namespace string, hfc *hfClientSet.Clientset) (vmList []hf.VirtualMachine, err error) {
	vmClaim, err := hfc.HobbyfarmV1().VirtualMachineClaims(Namespace).Get(context.TODO(), vmc, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			// object has been probably deleted.. no need to check
			return vmList, nil
		} else {
			return vmList, err
		}
	}

	dbr := vmClaim.Status.DynamicBindRequestId

	dbClaim, err := hfc.HobbyfarmV1().DynamicBindRequests(Namespace).Get(context.TODO(), dbr, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			return vmList, nil
		} else {
			return vmList, err
		}
	}

	for _, vm := range dbClaim.Status.VirtualMachineIds {
		vmInfo, err := hfc.HobbyfarmV1().VirtualMachines(Namespace).Get(context.TODO(), vm, metav1.GetOptions{})
		if err != nil {
			if apierrors.IsNotFound(err) {
				continue
			} else {
				return vmList, err
			}
		}
		vmList = append(vmList, *vmInfo)
	}

	return vmList, nil
}

func printReport(sessions []SessionDetails) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.Debug)
	fmt.Fprintln(w, "SESSION\t VMID\t STATUS\t PUBLICIP\t")
	for _, s := range sessions {
		for sessionID, vmList := range s.SessionVMMap {
			session := fmt.Sprintf("%s\t", sessionID)
			for _, vm := range vmList {
				output := fmt.Sprintf("%s %s\t %s\t %s\t", session, vm.Spec.Id, vm.Status.Status, vm.Status.PublicIP)
				fmt.Fprintln(w, output)
			}
		}
	}

	return w.Flush()
}

func printStats(sessions []SessionDetails) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.Debug)
	fmt.Fprintln(w, "TOTAL COUNT\t RUNNING VM's\t PENDING VM's\t")
	var total, running, pending int
	for _, session := range sessions {
		for _, vmList := range session.SessionVMMap {
			total = total + len(vmList)
			for _, vm := range vmList {
				if vm.Status.Status == "running" {
					running++
				} else {
					pending++
				}
			}
		}
	}

	fmt.Fprintf(w, "%d\t %d\t %d\t\n", total, running, pending)
	return w.Flush()
}
