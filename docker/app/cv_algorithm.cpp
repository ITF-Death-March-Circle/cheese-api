#include<opencv2/opencv.hpp>
#include"cv_algorithm.hpp"
constexpr int offset_width = 128;
constexpr int offset_height = 64;

//HSV�F�ϊ��p臒l
constexpr int B_MIN = 80;
constexpr int B_MAX = 200;
constexpr int G_MIN = 80;
constexpr int G_MAX = 220;
constexpr int R_MIN = 160;
constexpr int R_MAX = 255;
constexpr int  H_MAX = 30;
constexpr int  H_MIN = 0;
constexpr int  S_MAX = 150;
constexpr int  S_MIN = 30;
constexpr int  V_MAX = 255;
constexpr int  V_MIN = 60;

#define MIN_HSVCOLOR cv::Scalar(0, 60, 80)
#define MAX_HSVCOLOR cv::Scalar(10, 160, 240)

#define MIN_BKCOLOR cv::Scalar(0, 0, 0)
#define MAX_BKCOLOR cv::Scalar(180, 255, 40)


void extractFaceImage(std::string filepath, cv::UMat &clip_img) {
	cv::Mat _image = cv::imread(filepath);
	//cv::VideoCapture cap(0);
	cv::CascadeClassifier cascade; //�J�X�P�[�h���ފ�i�[�ꏊ
	cascade.load("haarcascade_frontalface_alt.xml"); //���ʊ��񂪓����Ă���J�X�P�[�h
	std::vector<cv::Rect> faces; //�֊s�����i�[�ꏊ
	cv::UMat dst;
	//cv::UMat clip_img;
	cv::UMat origin_dst;
	_image.copyTo(origin_dst);
	cv::convertScaleAbs(origin_dst, dst, 1.2, 80);

	//cv::imshow("diff", frame - before_frame);
	cascade.detectMultiScale(dst, faces, 1.1, 3, 0, cv::Size(50, 50));
	if (faces.size() == 0)return;
	cv::UMat work_src;

	cv::GaussianBlur(dst, dst, cv::Size(5, 5), 0);
	std::cout << faces.size() << std::endl;
	for (int i = 0; i < faces.size(); i++) //���o������̌�"faces.size()"�����[�v���s��
	{
		//�Z�[�t�N���b�s���O

		int w_start = (faces[i].x - offset_width) > 0 ? (faces[i].x - offset_width) : 0;
		int h_start = (faces[i].y - offset_height) > 0 ? (faces[i].y - offset_height) : 0;
		int w_end = (faces[i].x + faces[i].width + offset_width) < dst.cols ? faces[i].x + faces[i].width + offset_width : dst.cols;
		int h_end = (faces[i].y + faces[i].height + offset_height) < dst.rows ? faces[i].y + faces[i].height + offset_height : dst.rows;
		//rectangle(frame, cv::Point(faces[i].x, faces[i].y), cv::Point(faces[i].x + faces[i].width, faces[i].y + faces[i].height), cv::Scalar(0, 0, 255), 3); //���o�������ԐF��`�ň͂�
		//rectangle(frame, cv::Point(w_start, h_start), cv::Point(w_end, h_end), cv::Scalar(0, 0, 255), 3); //���o�������ԐF��`�ň͂�

		//�N���b�s���O�摜�̐���
		cv::UMat roi_img(origin_dst, cv::Rect(cv::Point(w_start, h_start), cv::Point(w_end, h_end)));
		//cv::imshow("clip" + std::to_string(i), roi_img);
		roi_img.copyTo(clip_img);
		roi_img.copyTo(work_src);

	}
	return;
	cv::UMat hsv_img;
	//cv::imshow("frame", frame);
	//cv::Mat diff;
	//cv::Mat gray_frame;
	//cv::Mat before_gray_frame;
	//cv::cvtColor(work_src, gray_frame, cv::COLOR_BGR2GRAY);
	//cv::cvtColor(before_frame, before_gray_frame, cv::COLOR_BGR2GRAY);
	//diff = gray_frame - before_gray_frame;
	//cv::threshold(diff, diff, 150, 255, cv::THRESH_BINARY|cv::THRESH_OTSU);
	//cv::imshow("diff2", diff);

	//cv::blur(org, dst, cv::Size(32, 32));
	//diff.copyTo(dst);
	//cv::GaussianBlur(org, dst, cv::Size(11, 11), 10, 10);
	if (!work_src.empty()) {
		cv::UMat grayImg;
		cv::UMat grayImg2;
		cv::UMat binaryImg;
		//diff.copyTo(grayImg);
/*		cv::Scalar s_min = cv::Scalar(B_MIN, G_MIN, R_MIN);
		cv::Scalar s_max = cv::Scalar(B_MAX, G_MAX, R_MAX);*/
		cv::Scalar s_min = cv::Scalar(H_MIN, S_MIN, V_MIN);
		cv::Scalar s_max = cv::Scalar(H_MAX, S_MAX, V_MAX);
		cv::cvtColor(work_src, hsv_img, cv::COLOR_BGR2HSV);
		// HSV�F��Ԃɂ����锧�F�̌��o
		cv::inRange(hsv_img, MIN_HSVCOLOR, MAX_HSVCOLOR, grayImg);
		/*cv::inRange(hsv_img, MIN_BKCOLOR, MAX_BKCOLOR, grayImg2);
		for (int y = 0; y < grayImg.rows; y++) {
			cv::Vec3b* src = grayImg.getMat(cv::ACCESS_RW).ptr<cv::Vec3b>(y);
			cv::Vec3b* diff = grayImg2.getMat(cv::ACCESS_RW).ptr<cv::Vec3b>(y);

			for (int x = 0; x < grayImg.cols; x++) {
				if (src[x][0] < 10 && diff[x][0] > 10) {
					src[x][0] = 255;
					src[x][1] = 255;
					src[x][2] = 255;

				}
			}
		}*/

		//work_src.copyTo(grayImg);
		//cv::cvtColor(work_src, grayImg, cv::COLOR_BGR2GRAY);
		//cv::equalizeHist(grayImg, grayImg);
		//cv::blur(grayImg, grayImg, cv::Size(8, 8));
			//�@���b�N�A�b�v�e�[�u���쐬
		int lut[256];

		//double gm = 1.0 / 2.0; //gamma�̓K���}�l
		//for (int i = 0; i < 256; i++)
		//{
		//	lut[i] = pow(1.0 * i / 255, gm) * 255;
		//}

		// �o�͉摜�ɓK�p
		//cv::LUT(grayImg, cv::Mat(cv::Size(256, 1), CV_8U, lut), grayImg);
		//cv::GaussianBlur(grayImg, grayImg, cv::Size(11, 11), 3, 3);
		//cv::Laplacian(grayImg, grayImg, CV_8U);

		//cv::medianBlur(grayImg, grayImg, 9);
		cv::Mat kernel = cv::getStructuringElement(cv::MORPH_RECT, cv::Size(5, 5));

		cv::morphologyEx(grayImg, grayImg2, cv::MORPH_GRADIENT, kernel, cv::Point(-1, -1), 4);

		cv::threshold(grayImg2, binaryImg, 150, 255, cv::THRESH_BINARY | cv::THRESH_OTSU);
		std::vector< std::vector< cv::Point > > contours;

		cv::findContours(binaryImg, contours, cv::RETR_EXTERNAL, cv::CHAIN_APPROX_TC89_L1);
		//cv::findContours(binaryImg, contours, cv::RETR_LIST, cv::CHAIN_APPROX_TC89_L1);
/*		for (auto contour = contours.begin(); contour != contours.end(); contour++) {
			cv::polylines(work_src, *contour, true, cv::Scalar(0, 255, 0), 2);
		}*/

		//�֊s�̐�
		int roiCnt = 0;

		//�֊s�̃J�E���g   
		int i = 0;

		for (auto contour = contours.begin(); contour != contours.end(); contour++) {
			std::vector< cv::Point > approx;

			//�֊s�𒼐��ߎ�����
			cv::approxPolyDP(cv::Mat(*contour), approx, 0.01 * cv::arcLength(*contour, true), true);

			// �ߎ��̖ʐς����ȏ�Ȃ�擾
			double area = cv::contourArea(approx);
			std::cout << area << std::endl;
			bool flag = false;
			if (area > 20000.0) {
				//�ň͂ޏꍇ            
				cv::polylines(dst, approx, true, cv::Scalar(255, 0, 0, 0), 2);
				std::stringstream sst;
				//sst << "area : " << area;
				//cv::putText(dst, sst.str(), approx[0], cv::FONT_HERSHEY_PLAIN, 1.0, cv::Scalar(0, 128, 0));

				//�֊s�ɗאڂ����`�̎擾
				cv::Rect brect = cv::boundingRect(cv::Mat(approx).reshape(2));
				//cv::drawContours(work_src, contours, i, CV_RGB(255, 0, 0), 1);
				//cv::UMat clip_img(work_src, brect);

				cv::UMat _clip_img(work_src, brect);
				_clip_img.copyTo(clip_img);
				//work_src.copyTo(clip_img);
				//cv::cvtColor(clip_img, clip_img, cv::COLOR_BGR2BGRA);

				//for (int j = 0; j < clip_img.rows; j++) {
				//	//std::cout << j << std::endl;
				//	flag = false;
				//	for (int i = 0; i < clip_img.cols; i++) {
				//		//std::cout << i << std::endl;
				//		if (clip_img.getMat(cv::ACCESS_RW).at<cv::Vec4b>(j, i)[2] == 255 && (clip_img.getMat(cv::ACCESS_RW).at<cv::Vec4b>(j, i)[1] + clip_img.getMat(cv::ACCESS_RW).at<cv::Vec4b>(j, i)[0]) == 0) {
				//			flag = !flag;
				//		}

				//		if (flag) {
				//			continue;
				//		}
				//		else {
				//			clip_img.getMat(cv::ACCESS_RW).at<cv::Vec4b>(j, i) = cv::Vec4b(0, 0, 0, 0);
				//		}
				//	}
				//}

				//cv::imshow("cliped", clip_img);
				//���͉摜�ɕ\������ꍇ

				roiCnt++;

				//�O�̂��ߗ֊s���J�E���g
				if (roiCnt == 99)
				{
					break;
				}
			}
			i++;
		}
	}
	//return clip_img;
}
